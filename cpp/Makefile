CXXFLAGS := -std=c++14 -ftime-report

.PHONY: all clean debug release

all: simple foldmap

debug: CXXFLAGS += -DDEBUG -g
debug: all

release: CXXFLAGS += -O2
release: all

clean:
	@rm -f simple foldmap