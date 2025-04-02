package endpoints

type FrontendPaths struct {
	ApiKeyForm   string
	ApiKeyVerify string
}

var FsPaths = FrontendPaths{
	ApiKeyForm:   "/form/api-key/",
	ApiKeyVerify: "/verify/api-key/",
}
