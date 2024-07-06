# Define the output names
OUT_HIDDEN=cleandl-daemon.exe
OUT_TERMINAL=CleanDL.exe
RC=cleandl-daemon.rc
RES=cleandl-daemon.res

# Build command for the hidden background task version
build-hidden:
	windres $(RC) -O coff -o $(RES)
	go build -ldflags="-H windowsgui -linkmode external -extldflags '-static -Wl,$(RES)'" -o $(OUT_HIDDEN)

# Build command for the terminal version 
build-terminal:
	go build -o $(OUT_TERMINAL)

# Default target to build both versions
all: build-hidden build-terminal