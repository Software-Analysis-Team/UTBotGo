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
.PHONY: update_package_config update_package_llvm_bitcode generate test update_answers $(foreach FN,$(FNS),update_answers_for_$(FN)) $(foreach FN,$(FNS),update_ta_config_for_$(FN)) $(foreach FN,$(FNS),update_Go_test_file_for_$(FN)) remove clean

$(DP)/name.txt $(DP)/import_path.txt $(DP)/go_files.txt $(DP)/pseudo_name.txt update_package_config &:
	go list -json | $(UUD)/package_config_parser Name $(DP)/name.txt ImportPath $(DP)/import_path.txt GoFiles $(DP)/go_files.txt
	$(UUD)/pseudo_name_getter $(DP)/name.txt $(DP)/import_path.txt $(DP)/pseudo_name.txt

$(DP)/file_templates.json:
	echo "{" > $@
	echo "    \"args\": \"$(foreach FN,%,$(DFFF))/args.json\"," >> $@
	echo "    \"answers\": \"$(foreach FN,%,$(DFFF))/answers.json\"," >> $@
	echo "    \"results\": \"$(foreach FN,%,$(DFFF))/results.json\"" >> $@
	echo "}" >> $@

$(DP)/llvm_bitcode.ll update_package_llvm_bitcode: $(DP)/go_files.txt
	$(UUD)/go2ll $(DP) $(DP)/llvm_bitcode.ll $<

$(foreach FN,$(FNS),$(DFFF)/klee_types.json) $(foreach FN,$(FNS),$(DFFF)/ta_types.json) &: $(DP)/go_files.txt
	$(UUD)/params_parser --go-files $$(cat $<) --klee $(foreach FN,%,$(DFFF))/klee_types.json --ta $(foreach FN,%,$(DFFF))/ta_types.json $(FNS)

define rules_for_function_name
$(DFFF)/llvm_bitcode.ll: $$(DP)/llvm_bitcode.ll $$(DP)/pseudo_name.txt
	$$(call mkdir,$$(@D))
	$$(UUD)/ll_modifier $$< $$@ $$$$(cat $$(DP)/pseudo_name.txt) $(FN)

$(DFFF)/klee-last: $(DFFF)/llvm_bitcode.ll
	klee --disable-verify --external-calls=all --entry-point=klee.main $$<

$(DFFF)/ta_config.json: $$(DP)/name.txt $$(DP)/file_templates.json $(DFFF)/ta_types.json $(DFFF)/args.json
	$$(UUD)/ta_config_creator $$^ > $$@

update_ta_config_for_$(FN): $$(DP)/name.txt $$(DP)/file_templates.json $(DFFF)/ta_types.json $(DFFF)/args.json update_answers_for_$(FN)
	$$(UUD)/ta_config_creator $$(DP)/name.txt $$(DP)/file_templates.json $(DFFF)/ta_types.json $(DFFF)/args.json $(DFFF)/answers.json > $(DFFF)/ta_config.json

$(DFFF)/args.json: $(DFFF)/klee_types.json $(DFFF)/klee-last
	ktest-tool --total-types $$< --to-json $$(wildcard $(DFFF)/klee-last/*.ktest) > $$@

$(FNTF): $(DFFF)/ta_config.json $$(UOFD)/test.go_template
	cd $$(UOFD); \
    $$(UUD)/template_applier -i $$(PWD)/$$< -c test.go_template $$(PWD)/$$@; \
    $$(UUD)/template_applier -i $$(PWD)/$$< $$(PWD)/$$@

update_Go_test_file_for_$(FN): $$(UOFD)/test.go_template update_ta_config_for_$(FN)
	cd $$(UOFD); \
    $$(UUD)/template_applier -i $$(PWD)/$(DFFF)/ta_config.json -c test.go_template $$(PWD)/$(FNTF); \
    $$(UUD)/template_applier -i $$(PWD)/$(DFFF)/ta_config.json $$(PWD)/$(FNTF)

$(DFFF)/results.json: | $(FNTF) $(DFFF)/args.json
	go test || true

update_answers_for_$(FN): $(DFFF)/results.json
	cp $$< $(DFFF)/answers.json

endef

$(eval $(foreach FN,$(FNS),$(rules_for_function_name)))

$(DP)/function_names.txt: $(UW)/config.yml
	rm -rf $(foreach FN,$(filter-out $(FNS),$(shell cat $@ 2> /dev/null)),$(DFFF)) $(foreach FN,$(filter-out $(FNS),$(shell cat $@ 2> /dev/null)),$(FNTF))
	echo $(FNS) > $@

generate: update_package_config update_package_llvm_bitcode $(foreach FN,$(FNS),$(FNTF)) $(foreach FN,$(FNS),$(DFFF)/args.json) $(DP)/function_names.txt

test: generate
	go test -v

update_answers: $(foreach FN,$(FNS),update_Go_test_file_for_$(FN))

clean:
	rm -rf $(UW) $(foreach FN,$(FNS),$(FNTF))
