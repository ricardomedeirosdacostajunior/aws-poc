package attachment

import (
	"aws-poc/internal/protocol"
	"fmt"
)

const filenameRoot = "disputes"

type (
	File struct {
		Key string
	}
	Service interface {
		Get(dispute *protocol.Dispute) (*protocol.Attachment, error)
		Save(chargeback *protocol.Chargeback) error
	}

	storage interface {
		list(cid string, bucket string, path string) ([]File, error)
		get(cid string, bucket string, key string) (*File, error)
	}

	repository interface {
		getUnsentFiles(*protocol.Dispute, []File) ([]File, error)
		save(chargeback *protocol.Chargeback) error
	}

	archiver interface {
		compact(cid string, files []File, strToRemove string) (*protocol.Attachment, error)
	}

	svc struct {
		storage
		archiver
		repository
	}
)

func (s svc) Get(dispute *protocol.Dispute) (*protocol.Attachment, error) {
	var (
		files          []File
		err            error
		unsentFiles    []File
		filesToCompact []File
	)
	path := fmt.Sprintf("%s/%d/%d", filenameRoot, dispute.AccountId, dispute.DisputeId)
	if files, err = s.list(dispute.Cid, dispute.OrgId, path); err != nil {
		return nil, err
	}
	if unsentFiles, err = s.getUnsentFiles(dispute, files); err != nil {
		return nil, err
	}

	for _, uf := range unsentFiles {
		var rf *File
		if rf, err = s.get(dispute.Cid, dispute.OrgId, uf.Key); err != nil {
			return nil, err
		}
		filesToCompact = append(filesToCompact, *rf)
	}

	return s.compact(dispute.Cid, filesToCompact, path)
}

func (s svc) Save(chargeback *protocol.Chargeback) error {
	return s.save(chargeback)
}
