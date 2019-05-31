package db

type nopLocker struct{}

func (nopLocker) Lock()    {}
func (nopLocker) Unlock()  {}
func (nopLocker) RLock()   {}
func (nopLocker) RUnlock() {}
