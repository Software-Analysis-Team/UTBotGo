CXXFLAGS = -O3 -Wall

.PHONY: all

all: bin/go2ll bin/ll_modifier bin/params_parser bin/ta_config_creator bin/package_configurator bin/package_config_parser bin/template_applier bin/config_parser

bin:
	mkdir $@

bin/ll_modifier_dir:
	mkdir $@

bin/go2ll: src/go2ll.sh | bin
	cp $< $@

bin/ll_modifier_dir/1: src/ll_modifier/1.cpp | bin/ll_modifier_dir
	$(CXX) $(CXXFLAGS) $$(llvm-config --cxxflags) -o $@ $< $$(llvm-config --ldflags --libs --system-libs)

bin/ll_modifier_dir/2: src/ll_modifier/2.go | bin/ll_modifier_dir
	go build -o $@ $<

bin/ll_modifier: src/ll_modifier.sh bin/ll_modifier_dir/1 bin/ll_modifier_dir/2
	cp $< $@

bin/params_parser: $(wildcard src/params_parser/*.go) | bin
	cd src/params_parser; go build -o ../../$@ .

bin/ta_config_creator: src/ta_config_creator.py | bin
	cp $< $@

bin/package_configurator: src/package_configurator.sh | bin
	cp $< $@

bin/package_config_parser: src/package_config_parser.py | bin
	cp $< $@

bin/template_applier: $(wildcard src/template_applier/*.go) | bin
	cd src/template_applier; go build -o ../../$@ .

bin/config_parser: src/config_parser.py | bin
	cp $< $@

clean:
	rm -r bin
