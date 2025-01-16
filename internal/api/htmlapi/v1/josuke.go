package v1

import (
	"context"
	"log/slog"
	"monkeydioude/grig/internal/html/blocks"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/internal/model"
	htmlinput "monkeydioude/grig/internal/service/html_input"
	"monkeydioude/grig/internal/service/utils"
	"monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/html/elements"
	pkgModel "monkeydioude/grig/pkg/model"
	"net/http"
	"strconv"
)

func (h Handler) JosukeList(w http.ResponseWriter, r *http.Request, _ *slog.Logger, nav elements.Nav) error {
	layout := layouts.Main(nav, pages.JosukeList(&h.Layout.ServerConfig))
	return layout.Render(context.Background(), w)
}

func (h Handler) JosukeHookBlock(w http.ResponseWriter, r *http.Request, _ *slog.Logger) error {
	indexStr := r.URL.Query().Get("index")
	index := 0
	if indexStr != "" {
		it, err := strconv.Atoi(indexStr)
		if err != nil {
			return errors.Wrap(err, "JosukeHookBlock")
		}
		index = it
	}

	return blocks.JosukeHook(index, pages.GetHookInputName, model.Hook{}).Render(context.Background(), w)
}

func (h Handler) JosukeDeploymentBlock(w http.ResponseWriter, r *http.Request, _ *slog.Logger) error {
	indexStr := r.URL.Query().Get("index")
	index := 0
	if indexStr != "" {
		it, err := strconv.Atoi(indexStr)
		if err != nil {
			return errors.Wrap(err, "JosukeDeploymentBlock")
		}
		index = it
	}

	dep := model.Deployment{
		Indexer: model.Indexer{Index: index},
	}

	dep.FillBaseData()
	return blocks.JosukeDeployment(pages.GetInputNameWithKey, utils.GetInputName, dep).Render(context.Background(), w)
}

func (h Handler) JosukeBranchBlock(w http.ResponseWriter, r *http.Request, _ *slog.Logger) error {
	indexStr := r.URL.Query().Get("index")
	parentNameStr := r.URL.Query().Get("parent_name")

	branch := model.NewBranch(0)
	josukeTree(branch, indexStr, parentNameStr)
	htmlinput.FillIndexesFromHTMLInput(branch.GetParent(), parentNameStr)
	return blocks.JosukeBranch(utils.GetInputName, *branch).Render(context.Background(), w)
}

func (h Handler) JosukeActionBlock(w http.ResponseWriter, r *http.Request, _ *slog.Logger) error {
	indexStr := r.URL.Query().Get("index")
	parentNameStr := r.URL.Query().Get("parent_name")
	action := model.NewAction(0)
	josukeTree(action, indexStr, parentNameStr)
	htmlinput.FillIndexesFromHTMLInput(action.GetParent(), parentNameStr)
	return blocks.JosukeAction(utils.GetInputName, *action).Render(context.Background(), w)
}

func (h Handler) JosukeCommandBlock(w http.ResponseWriter, r *http.Request, _ *slog.Logger) error {
	indexStr := r.URL.Query().Get("index")
	parentNameStr := r.URL.Query().Get("parent_name")
	cmd := model.NewCommand(0)
	josukeTree(cmd, indexStr, parentNameStr)
	htmlinput.FillIndexesFromHTMLInput(cmd.GetParent(), parentNameStr)
	return blocks.JosukeCmd(utils.GetInputName, *cmd).Render(context.Background(), w)
}

func josukeTree(hp pkgModel.IndexBuilder, indexStr string, parentNameStr string) error {
	index := 0
	if indexStr != "" {
		it, err := strconv.Atoi(indexStr)
		if err != nil {
			return errors.Wrap(err, "JosukeBranchAction")
		}
		index = it
	}

	hp.SetIndex(index)
	htmlinput.FillIndexesFromHTMLInput(hp.GetParent(), parentNameStr)
	return nil
}
