package command

type Handler struct {
	fileTransfer FileTransferService
}

func NewHandler(
	fileTransfer FileTransferService,
) *Handler {
	return &Handler{
		fileTransfer: fileTransfer,
	}
}
