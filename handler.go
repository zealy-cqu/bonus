package bonus

//import (
//	"sync"
//)
//
//var hdEntry = handlerEntry{
//	CSID_ID_HeartbeatRsp: handleHeartbeatRsp,
//	CSID_ID_WinningNtf: handleWinningNtf,
//}
//
//
//type handleFn func(msg *message) error
//type handlerEntry map[CSID]handleFn
//
//
//type handler struct {
//	hdEntry handlerEntry
//
//	wg *sync.WaitGroup
//	closeCh chan struct{}
//
//	ss *session
//	readCh chan *readData
//	readErrCh chan error
//
//	handleErrCh chan *handleErr
//}
//
//
//type handleErr = readData
//
//
////func newHandler() *handler {
////	return nil
////}
//
//func (h *handler) handle()  {
//	h.wg.Add(1)
//	defer h.wg.Done()
//
//	go h.ss.read(h.readCh)
//
//	for {
//		select {
//		case <- h.closeCh:
//			return
//		case in := <- h.readCh:
//			if in.err != nil {
//				h.readErrCh <- in.err
//			}
//			if err := h.hdEntry[in.msg.id](in.msg); err !=  nil {
//				h.handleErrCh <- &handleErr{in.msg, err}
//			}
//		}
//	}
//}
//
//// Note that close will stop handle but won't close the underlying session.
//func (h *handler) close() {
//	close(h.closeCh)
//}
//
//
//func handleHeartbeatRsp(msg *message) error {
//	return nil
//}
//
//
//func handleWinningNtf(msg *message) error {
//	return nil
//}
