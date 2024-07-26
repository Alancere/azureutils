package typespecgo

// https://github.com/Azure/autorest.go/blob/main/packages/typespec-go/src/lib.ts#GoEmitterOptions
/*
@azure-tools/typespec-go option
*/
type GoEmitterOptions struct {
	AzcoreVersion           string `yaml:"azcore-version"`
	DisallowUnknownFields   bool   `yaml:"disallow-unknown-fields"`
	FilePrefix              string `yaml:"file-prefix"`
	GenerateFake            bool   `yaml:"generate-fakes"`
	InjectSpanc             bool   `yaml:"inject-spans"`
	Module                  string `yaml:"module"`
	ModuleVersion           string `yaml:"module-version"`
	RawJsonAsBytes          bool   `yaml:"rawjson-as-bytes"`
	SliceElementsByVal      bool   `yaml:"slice-elements-byval"`
	SingleClient            bool   `yaml:"single-client"`
	Stutter                 string `yaml:"stutter"`
	FixConstStuttering      bool   `yaml:"fix-const-stuttering"`
	RemoveUnreferencedTypes bool   `yaml:"remove-unreferenced-types"`
}

type EmitOption struct {
	ServiceDir       string `yaml:"service-dir"`
	PackageDir       string `yaml:"package-dir"`
	EmitterOutputDir string `yaml:"emitter-output-dir"`
}

type GoOption struct {
	EmitOption

	GoEmitterOptions
}
