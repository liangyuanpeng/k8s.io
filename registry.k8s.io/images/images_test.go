package images

import (
	"log"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func TestPullImage(t *testing.T) {
	// 仓库路径
	repoPath := "/home/runner/work/lanactions/lanactions/k8s.io"

	// 分支名称
	baseBranch := "main"
	targetBranch := "patch-1"

	// 打开仓库
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	// 获取两个分支的Commit对象
	baseRef, err := repo.Reference(plumbing.NewBranchReferenceName(baseBranch), true)
	if err != nil {
		log.Fatal(err)
	}
	baseCommit, err := repo.CommitObject(baseRef.Hash())
	if err != nil {
		log.Fatal(err)
	}
	targetRef, err := repo.Reference(plumbing.NewBranchReferenceName(targetBranch), true)
	if err != nil {
		log.Fatal(err)
	}
	targetCommit, err := repo.CommitObject(targetRef.Hash())
	if err != nil {
		log.Fatal(err)
	}
	patch, err := baseCommit.Patch(targetCommit)
	log.Println("patch.files:", len(patch.FilePatches()))

	// 获取两个Commit之间的文件差异
	// filePatches, err := object.Diff(baseCommit, targetCommit)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	for _, v := range patch.FilePatches() {
		from, to := v.Files()
		// log.Println("from:", from == nil, to == nil)
		fromPath := ""
		if from != nil {
			fromPath = from.Path()
		}
		toPath := ""
		if to != nil {
			toPath = to.Path()
		}
		log.Printf("filepatch:%s|%s \n", fromPath, toPath)
		if !strings.HasPrefix(toPath, "registry.k8s.io/images") {
			continue
		}
		if !strings.HasSuffix(toPath, ".yaml") {
			continue
		}
		log.Printf("filepatch2:%s|%s \n", fromPath, toPath)
	}

	// 遍历差异的文件
	// patch.FilePatches().ForEach(func(filePatch *object.FilePatch) error {
	// 	// 输出文件路径
	// 	fmt.Println("File:", filePatch.Name)

	// 	// 输出文件的差异内容
	// 	content, err := filePatch.String()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Println(content)

	// 	return nil
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
