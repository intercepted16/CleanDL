# Define the output names
OUT_HIDDEN=background_task.exe
OUT_TERMINAL=CleanDL.exe
RC=background_task.rc

# Build command for the hidden background task version
build-hidden:
	windres $(RC) -O coff -o background_task.res
	go build -ldflags="-H windowsgui -linkmode external -extldflags '-static -Wl,background_task.res'" -o $(OUT_HIDDEN)

# Build command for the terminal version 
build-terminal:
	go build -o $(OUT_TERMINAL)

# Default target to build both versions
all: build-hidden build-terminal