package core

// Rebuild recomputes the item list from the current directory and config, then rebuilds the
// table columns and rows. Call after any state change that affects the listing: sort, group,
// filter, hidden toggle, navigation into a directory. ctx.Current must be non-nil.
func Rebuild(ctx *Context) {
	items := ctx.Current.ComputeItems(ctx.Config)
	ctx.Items = items
	buildColumns(ctx)
	buildRows(ctx, items)
}

// RestoreCursorByUID positions the cursor on the row whose UID matches the current directory's
// lastSelectedUID, when one is set. The set call happens just before navigation so the cursor
// lands back on the same logical row the user came from. ctx.Current must be non-nil.
func RestoreCursorByUID(ctx *Context) {
	uid := ctx.Current.LastSelectedUID()
	if uid == 0 {
		return
	}
	for i, item := range ctx.Items {
		if item.UID() == uid {
			ctx.Table.SetCursor(i)
			return
		}
	}
}
