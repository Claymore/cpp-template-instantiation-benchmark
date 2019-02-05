package main

import (
	"fmt"
	"os"
	"text/template"
)

const (
	structsHeader = `#pragma once
`
	structTemplateString = `
struct {{.Name}} {
	int x = 0;
};
`
	simpleCppHeader = `#include <algorithm>
#include <vector>

#include "structs.h"

int main() {
	return 0;
}
`
	simpleLoopString = `
int TestFunc{{.Name}}(const std::vector<int>& src, std::vector<{{.Name}}>& dest) {
	dest.clear();
	dest.reserve(src.size());
	for (auto it = src.begin(); it != src.end(); ++it) {
		{{.Name}} tmp;
		tmp.x = *it;
		dest.emplace_back(std::move(tmp));
	}
	return 0;
}
`
	foldmapHppHeader = `#pragma once
#include <algorithm>
#include <vector>

template<class Fold, class Map, class InputIt, class Dest>
void FoldMap(const Fold& accumulateIn, const Map& transform, InputIt first, InputIt last, Dest* destination) {
    std::for_each(first, last, [&](auto& item) {
        accumulateIn(destination, transform(item));
    });
}

template<class T, class ConvertorFunction>
void ConvertToVector(const std::vector<int>& a, std::vector<T>* result, ConvertorFunction f) {
    result->reserve(a.size());
    auto convertItem = [f](auto& t) {
        T tmp;
        f(t, tmp);
        return tmp;
    };
    auto appendToDest = [](auto* result, const auto& tmp) {
        result->emplace_back(std::move(tmp));
    };
    FoldMap(appendToDest, convertItem, a.begin(), a.end(), result);
}
`
	foldmapCppHeader = `#include "foldmap.hpp"

#include <algorithm>
#include <vector>

#include "structs.h"

int main() {
	return 0;
}
`
	foldMapString = `
int TestFunc{{.Name}}(const std::vector<int>& src, std::vector<{{.Name}}>& dest) {
	ConvertToVector(src, &dest, [](auto& t, auto& tmp) {
		tmp.x = t;
	});
	return 0;
}
`
)

func main() {
	const n = 300

	type codeExample struct {
		Name string
	}

	structTemplate := template.Must(template.New("struct").Parse(structTemplateString))
	foldMapTemplate := template.Must(template.New("foldMap").Parse(foldMapString))
	simpleLoopTemplate := template.Must(template.New("simpleLoop").Parse(simpleLoopString))

	structsH, _ := os.Create("cpp/structs.h")
	defer structsH.Close()
	structsH.Write([]byte(structsHeader))

	simpleCpp, _ := os.Create("cpp/simple.cpp")
	defer simpleCpp.Close()
	simpleCpp.Write([]byte(simpleCppHeader))

	foldmapHpp, _ := os.Create("cpp/foldmap.hpp")
	defer foldmapHpp.Close()
	foldmapHpp.Write([]byte(foldmapHppHeader))

	foldmapCpp, _ := os.Create("cpp/foldmap.cpp")
	defer foldmapCpp.Close()
	foldmapCpp.Write([]byte(foldmapCppHeader))

	for i := 0; i < n; i++ {
		example := codeExample{Name: fmt.Sprintf("A%d", i+1)}
		structTemplate.Execute(structsH, example)
		simpleLoopTemplate.Execute(simpleCpp, example)
		foldMapTemplate.Execute(foldmapCpp, example)
	}
}
