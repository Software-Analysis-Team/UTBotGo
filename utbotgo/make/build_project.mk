targets := $(shell ls -A $(TEMPLATES_DIR))

# This is function for creating templated rule. Recipe is reading content from
# file-template, replacing all variables to corresponding values from makefile
# `config.mk` by using `$(SUB_PROGRAM_DIR)/$(SUB_PROGRAM_NAME)` and writing
# result to file-target
define rule_for_target
$(1): $(TEMPLATES_DIR)/$(1) $(SUB_PROGRAM_DIR)/$(SUB_PROGRAM_NAME)
	cat $$< | $(SUB_PROGRAM_DIR)/$(SUB_PROGRAM_NAME) $(CONFIG_FILE) > $$@

endef

$(eval $(foreach target,$(targets),$(call rule_for_target,$(target))))

build_project: $(targets)
