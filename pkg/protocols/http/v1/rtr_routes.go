package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *HttpEndpoints) AddRTRSpecificEndpoints(rg *gin.RouterGroup) {
	auth := rg.Group("/rtr")
	auth.GET("/code-validation", h.rtrCodeValidationEndpoint)
}

func (h *HttpEndpoints) rtrCodeValidationEndpoint(c *gin.Context) {
	code, codeExists := c.GetQuery("code")

	if !codeExists || len(code) < 4 {
		time.Sleep(time.Second * 2)
		c.JSON(http.StatusBadRequest, gin.H{"error": "code error"})
		return
	}

	for k, codes := range codeLists {
		for _, ec := range codes {
			if ec == code {
				c.JSON(http.StatusOK, gin.H{
					"group": k,
					"code":  code,
				})
				return
			}
		}
	}

	time.Sleep(time.Second * 2)
	c.JSON(http.StatusBadRequest, gin.H{"error": "code error"})
}

var codeLists map[string][]string = map[string][]string{
	"KG0829": {
		"K101",
		"K102",
		"K103",
		"K104",
		"K105",
		"K106",
		"K107",
		"K108",
		"K109",
		"K110",
		"K111",
		"K112",
		"K113",
		"K114",
		"K115",
		"K116",
		"K117",
		"K118",
		"K119",
		"K120",
		"K121",
		"K122",
		"K123",
		"K124",
		"K125",
		"K126",
		"K127",
		"K128",
		"K129",
		"K130",
	},
	"KG0912": {
		"K201",
		"K203",
		"K202",
		"K204",
		"K205",
		"K206",
		"K207",
		"K208",
		"K209",
		"K210",
		"K211",
		"K212",
		"K213",
		"K214",
		"K215",
		"K216",
		"K217",
		"K218",
		"K219",
		"K220",
		"K221",
		"K222",
		"K223",
		"K224",
		"K225",
		"K226",
		"K227",
		"K228",
		"K229",
		"K230",
	},
	"EG": {
		"E001",
		"E002",
		"E003",
		"E004",
		"E005",
		"E006",
		"E007",
		"E008",
		"E009",
		"E010",
		"E011",
		"E012",
		"E013",
		"E014",
		"E015",
		"E016",
		"E017",
		"E018",
		"E019",
		"E020",
		"E021",
		"E022",
		"E023",
		"E024",
		"E025",
		"E026",
		"E027",
		"E028",
		"E029",
		"E030",
		"E031",
		"E032",
		"E033",
		"E034",
		"E035",
		"E036",
		"E037",
		"E038",
		"E039",
		"E040",
		"E041",
		"E042",
		"E043",
		"E044",
		"E045",
		"E046",
		"E047",
		"E048",
		"E049",
		"E050",
		"E051",
		"E052",
		"E053",
		"E054",
		"E055",
		"E056",
		"E057",
		"E058",
		"E059",
		"E060",
		"E061",
		"E062",
		"E063",
		"E064",
		"E065",
		"E066",
		"E067",
		"E068",
		"E069",
		"E070",
		"E071",
		"E072",
		"E073",
		"E074",
		"E075",
		"E076",
		"E077",
		"E078",
		"E079",
		"E080",
		"E081",
		"E082",
		"E083",
		"E084",
		"E085",
		"E086",
		"E087",
		"E088",
		"E089",
		"E090",
		"E091",
		"E092",
		"E093",
		"E094",
		"E095",
		"E096",
		"E097",
		"E098",
		"E099",
	},
}
