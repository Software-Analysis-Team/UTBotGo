# Functions for using marks
set_mark = mkdir -p $(MARK_DIR); touch $(MARK_DIR)/$(1)
reset_mark = rm -f $(MARK_DIR)/$(1)
mark_is_set = [ -f $(MARK_DIR)/$(1) ]
mark_is_not_set = [ ! -f $(MARK_DIR)/$(1) ]
