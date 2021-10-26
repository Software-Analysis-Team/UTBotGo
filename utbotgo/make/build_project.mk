targets := $(shell ls -A $(TEMPLATES_DIR))

define rule_for_target
$(1): $(TEMPLATES_DIR)/$(1)
	cat $$^ | $(SUB_PROGRAM_DIR)/$(SUB_PROGRAM_NAME) $(MAKE_DIR)/config.mk > $$@

endef

$(eval $(foreach target,$(targets),$(call rule_for_target,$(target))))

build_project: $(targets)
