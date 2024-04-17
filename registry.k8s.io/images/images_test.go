package images

import (
	"context"
	"log"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"
	"sigs.k8s.io/promo-tools/v4/image"
)

func TestPullImage(t *testing.T) {
	// go test -timeout 30s -run ^TestPullImage$ images -v
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
		log.Fatal("fetch "+baseBranch+" failed!", err)
	}
	baseCommit, err := repo.CommitObject(baseRef.Hash())
	if err != nil {
		log.Fatal(err)
	}
	targetRef, err := repo.Reference(plumbing.NewBranchReferenceName(targetBranch), true)
	if err != nil {
		log.Fatal("fetch "+targetBranch+" failed!", err)
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

		fromManifestList := &image.ManifestList{}
		toManifestList := &image.ManifestList{}

		if fromPath != "" {
			fromManifestList, err = image.NewManifestListFromFile(filepath.Join(repoPath + "/" + fromPath))
			if err != nil {
				panic(err)
			}
		}

		toManifestList, err = image.NewManifestListFromFile(filepath.Join(repoPath + "/" + toPath))
		if err != nil {
			panic(err)
		}

		newManifestList := diffNewManifestList(fromManifestList, toManifestList)
		for _, imageData := range *newManifestList {
			log.Println("imageData:", imageData.Name)
		}

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

func diffNewManifestList(from, to *image.ManifestList) *image.ManifestList {
	return &image.ManifestList{}
}

func TestImages(t *testing.T) {
	imgfilepath := "/home/runner/work/lanactions/lanactions/k8s.io/registry.k8s.io/images"
	newManifestList, err := image.NewManifestListFromFile(filepath.Join(imgfilepath))
	if err != nil {
		panic(err)
	}
	reg, err := remote.NewRegistry("gcr.io")
	if err != nil {
		panic(err)
	}
	total := 0
	checkDigestFailedTags := []string{}
	for _, imageData := range *newManifestList {
		src, err := reg.Repository(context.TODO(), "k8s-staging-dns/"+imageData.Name)
		if err != nil {
			panic(err)
		}
		for digest, tags := range imageData.DMap {
			for _, tag := range tags {
				if tag != "1.23.0" {
					continue
				}
				log.Println("checking image for:", "gcr.io/k8s-staging-dns/"+imageData.Name+":"+tag)
				total++
				dst := memory.New()
				desc, err := oras.Copy(context.TODO(), src, tag, dst, tag, oras.DefaultCopyOptions)
				if err != nil {
					panic(err)
				}
				if desc.Digest.String() != digest {
					checkDigestFailedTags = append(checkDigestFailedTags, "gcr.io/k8s-staging-dns/"+imageData.Name+":"+tag)
				}
			}
		}
	}
	log.Printf("check total:%d, failed:%d \n", total, len(checkDigestFailedTags))

}
