package token

const (
	ILLEGAL = iota
	EOF
	// taken from the dockerfile reference
	ADD              //Add local or remote files and directories.
	ARG              //Use build-time variables.
	CMD              //Specify default commands.
	COPY             //Copy files and directories.
	ENTRYPOINT       //Specify default executable.
	ENV              //Set environment variables.
	EXPOSE           //Describe which ports your application is listening on.
	FROM             //Create a new build stage from a base image.
	HEALTHCHECK      //Check a container's health on startup.
	LABEL            //Add metadata to an image.
	MAINTAINER       //Specify the author of an image.
	ONBUILD          //Specify instructions for when the image is used in a build.
	RUN              //Execute build commands.
	SHELL            //Set the default shell of an image.
	STOPSIGNAL       //Specify the system call signal for exiting a container.
	USER             //Set user and group ID.
	VOLUME           //Create volume mounts.
	WORKDIR          //Change working directory.
	COMMENT          //Comment line
	PARSER_DIRECTIVE //Comment line with parser directive data
)

var TokenLookupTable = map[string]int{
	"ILLEGAL":     ILLEGAL, // Can
	"EOF":         EOF,
	"ADD":         ADD,
	"ARG":         ARG,
	"CMD":         CMD,
	"COPY":        COPY,
	"ENTRYPOINT":  ENTRYPOINT,
	"ENV":         ENV,
	"EXPOSE":      EXPOSE,
	"FROM":        FROM,
	"HEALTHCHECK": HEALTHCHECK,
	"LABEL":       LABEL,
	"MAINTAINER":  MAINTAINER,
	"ONBUILD":     ONBUILD,
	"RUN":         RUN,
	"SHELL":       SHELL,
	"STOPSIGNAL":  STOPSIGNAL,
	"USER":        USER,
	"VOLUME":      VOLUME,
	"WORKDIR":     WORKDIR,
}

type Token struct {
	Kind          int
	Params        map[string]string
	Content       string
	InlineComment string
}
