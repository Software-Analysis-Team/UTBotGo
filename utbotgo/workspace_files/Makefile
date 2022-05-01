SHELL = /usr/bin/env bash

UUD = /utbotgo/utils/bin
UW = utbotgo
UOFD = /utbotgo/utils/other_files
PWD := $(shell pwd)
DP = $(UW)/package
DFFF = $(UW)/functions/$(FN)
FNTF = utbotgo_$(FN)_test.go

mkdir = mkdir -p $(1)

# Rules
.PHONY: update_package_config update_package_llvm_bitcode generate test update_answers $(foreach FN,$(FNS),update_answers_for_$(FN)) remove clean

$(DP)/config.json: update_package_config

update_package_config:
	$(call mkdir,$(DP))
	$(UUD)/package_configurator $(DP)/config.json

$(DP)/real_name.txt: $(DP)/config.json
	$(UUD)/package_config_parser Name < $< > $@

$(DP)/name.txt: $(DP)/config.json $(DP)/real_name.txt
	name=$$(cat $(DP)/real_name.txt); \
    if [[ $$name != 'main' ]]; then \
        name=$$($(UUD)/package_config_parser ImportPath < $<); \
    fi; \
    echo "$$name" > $@

$(DP)/go_files.txt: $(DP)/config.json
	$(UUD)/package_config_parser GoFiles < $< > $@

$(DP)/file_templates.json:
	echo "{" > $@
	echo "    \"args\": \"$(foreach FN,%,$(DFFF))/args.json\"," >> $@
	echo "    \"answers\": \"$(foreach FN,%,$(DFFF))/answers.json\"," >> $@
	echo "    \"results\": \"$(foreach FN,%,$(DFFF))/results.json\"" >> $@
	echo "}" >> $@

$(DP)/llvm_bitcode.ll: update_package_llvm_bitcode

update_package_llvm_bitcode: $(DP)/go_files.txt
	$(call mkdir,$(DP))
	$(UUD)/go2ll $(DP) $(DP)/llvm_bitcode.ll $<

$(foreach FN,$(FNS),$(DFFF)/klee_types.json) $(foreach FN,$(FNS),$(DFFF)/ta_types.json) &:: $(DP)/go_files.txt
	$(UUD)/params_parser --go-files $$(cat $<) --klee $(foreach FN,%,$(DFFF))/klee_types.json --ta $(foreach FN,%,$(DFFF))/ta_types.json $(FNS)

define rules_for_function_name
$(DFFF)/llvm_bitcode.ll: $$(DP)/llvm_bitcode.ll $$(DP)/name.txt
	$$(call mkdir,$$(@D))
	$$(UUD)/ll_modifier $$< $$@ $$$$(cat $$(DP)/name.txt) $(FN)

$(DFFF)/klee-last: $(DFFF)/llvm_bitcode.ll
	klee --disable-verify --external-calls=all --entry-point=klee.main $$<

$(DFFF)/ta_config.json: $$(DP)/real_name.txt $$(DP)/file_templates.json $(DFFF)/ta_types.json
	$$(UUD)/ta_config_creator $$^ > $$@

$(DFFF)/args.json: $(DFFF)/klee_types.json $(DFFF)/klee-last
	ktest-tool --total-types $$< --to-json $$(wildcard $(DFFF)/klee-last/*.ktest) > $$@

$(FNTF): $(DFFF)/ta_config.json $$(UOFD)/test.go_template
	cd $$(UOFD); $$(UUD)/template_applier -i $$(PWD)/$$< -c test.go_template $$(PWD)/$$@

$(DFFF)/results.json: $(FNTF) $(DFFF)/args.json
	go test || true

update_answers_for_$(FN): $(DFFF)/results.json
	cp $$< $(DFFF)/answers.json

endef

$(eval $(foreach FN,$(FNS),$(rules_for_function_name)))

$(DP)/function_names.txt: $(UW)/config.yml
	rm -rf $(foreach FN,$(filter-out $(FNS),$(shell cat $@ 2> /dev/null)),$(DFFF)) $(foreach FN,$(filter-out $(FNS),$(shell cat $@ 2> /dev/null)),$(FNTF))
	echo $(FNS) > $@

generate: $(foreach FN,$(FNS),$(FNTF)) $(foreach FN,$(FNS),$(DFFF)/args.json) $(DP)/function_names.txt

test: generate
	go test -v || true

update_answers: $(foreach FN,$(FNS),update_answers_for_$(FN))

clean:
	rm -rf $(UW) $(foreach FN,$(FNS),$(FNTF))
