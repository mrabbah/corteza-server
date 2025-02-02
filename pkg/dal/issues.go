package dal

type (
	issue struct {
		kind issueKind
		err  error
	}
	issueSet []issue

	issueHelper struct {
		// these two will be used to help clear out unneeded errors
		connections []uint64
		models      []uint64

		connectionIssues dalIssueIndex
		modelIssues      dalIssueIndex
	}

	issueKind     string
	dalIssueIndex map[uint64]issueSet
)

const (
	connectionIssue issueKind = "connection"
	modelIssue      issueKind = "model"
)

func newIssueHelper() *issueHelper {
	return &issueHelper{
		connectionIssues: make(dalIssueIndex),
		modelIssues:      make(dalIssueIndex),
	}
}

func makeIssue(kind issueKind, err error) issue {
	return issue{
		kind: kind,
		err:  err,
	}
}

func (svc *service) SearchConnectionIssues(connectionID uint64) (out []error) {
	for _, issue := range svc.connectionIssues[connectionID] {
		out = append(out, issue.err)
	}

	return
}

func (svc *service) SearchModelIssues(connectionID, resourceID uint64) (out []error) {
	// @todo index by connection as well
	// if _, ok := svc.modelIssues[connectionID]; !ok {
	// 	return
	// }

	for _, issue := range svc.modelIssues[resourceID] {
		out = append(out, issue.err)
	}

	return
}

func (svc *service) hasConnectionIssues(connectionID uint64) bool {
	return len(svc.SearchConnectionIssues(connectionID)) > 0
}

func (svc *service) hasModelIssues(connectionID, modelID uint64) bool {
	return len(svc.SearchModelIssues(connectionID, modelID)) > 0
}

func (svc *service) updateIssues(issues *issueHelper) {
	for _, connectionID := range issues.connections {
		delete(svc.connectionIssues, connectionID)
	}
	for connectionID, issues := range issues.connectionIssues {
		svc.connectionIssues[connectionID] = issues
	}

	for _, modelID := range issues.models {
		delete(svc.modelIssues, modelID)
	}
	for modelID, issues := range issues.modelIssues {
		svc.modelIssues[modelID] = issues
	}
}

func (svc *service) clearModelIssues() {
	svc.modelIssues = make(dalIssueIndex)
}

func (rd *issueHelper) addConnection(connectionID uint64) *issueHelper {
	rd.connections = append(rd.connections, connectionID)
	return rd
}

func (rd *issueHelper) addModel(modelID uint64) *issueHelper {
	rd.models = append(rd.models, modelID)
	return rd
}

func (rd *issueHelper) addConnectionIssue(connectionID uint64, err error) {
	rd.connectionIssues[connectionID] = append(rd.connectionIssues[connectionID], makeIssue(connectionIssue, err))
}

func (rd *issueHelper) addModelIssue(resourceID uint64, err error) {
	rd.modelIssues[resourceID] = append(rd.modelIssues[resourceID], makeIssue(modelIssue, err))
}

func (a *issueHelper) mergeWith(b *issueHelper) {
	if b == nil {
		return
	}

	for connectionID, issues := range b.connectionIssues {
		a.connectionIssues[connectionID] = append(a.connectionIssues[connectionID], issues...)
	}

	for modelID, issues := range b.modelIssues {
		a.modelIssues[modelID] = append(a.modelIssues[modelID], issues...)
	}
}

// Op check utils

func (svc *service) canOpData(connectionID, modelID uint64) (err error) {
	if svc.hasConnectionIssues(connectionID) {
		return errRecordOpProblematicConnection(connectionID)
	}
	if svc.hasModelIssues(connectionID, modelID) {
		return errRecordOpProblematicModel(modelID)
	}

	return nil
}
