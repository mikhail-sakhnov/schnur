package targets

var predifinedTypes = commandsByTargets{
	"vagrant_test": CommandsList{
		"sleep 5",
		"uname -a",
	},
}
