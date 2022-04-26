package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

chart: schema.#resource & {
	parents: [
		{handle: "namespace"},
	]

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
		}
	}

	locale: {
		extended: true

		keys: {
			configReportsYaxis: {
				path: ["reports", {part: "reportID", var: true}, "yAxis", "label"]
				customHandler: true
			}
		}
	}
}
