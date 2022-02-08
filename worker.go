package dilema

func (di *dicon) goOperationIndexProvider() {
	var operationIndex uint64 = 0
	for {
		di.operationIndexCh <- operationIndex
		operationIndex++
	}
}

func (di *dicon) goQueueWriter() {
	for {
		select {
		case event := <-di.queueCh:
			di.pushEventBack(event)
		case <-di.exitCh:
			return
		}
	}
}

func (di *dicon) goQueueReader() {
	workersCount := 1
	go di.goDiconWorker()
	for {
		if qLen := di.queueLen(); qLen > 0 {
			if qLen > 2*workersCount {
				workersCount++
				go di.goDiconWorker()
			}
			if qLen < workersCount {
				di.exitCh <- struct{}{}
			}
			di.operationStartCh <- di.popEvent()
		}
	}
}

func (di *dicon) goDiconWorker() {
	for {
		select {
		case startEvent := <-di.operationStartCh:
			operationIndex := startEvent.operationIndex
			switch startEvent.oType {
			case registerTemporalOperation:
				event := startEvent.event.(registerTemporalStartEvent)
				err := di.registerTemporal(event.alias, event.serviceInit)
				di.registerEndCh <- operationEndEvent{
					operationIndex,
					registerEndEvent{err: err},
				}
			case registerSingleToneOperation:
				event := startEvent.event.(registerSingleToneStartEvent)
				err := di.registerSingleTone(event.alias, event.serviceInit, event.args...)
				di.registerEndCh <- operationEndEvent{
					operationIndex,
					registerEndEvent{err: err},
				}
			case registerFewOperation:
				event := startEvent.event.(registerFewStartEvent)
				err := di.registerFew(event.servicesInit, event.args...)
				di.registerEndCh <- operationEndEvent{
					operationIndex,
					registerEndEvent{err: err},
				}
			case getSingleToneOperation:
				event := startEvent.event.(getSingleToneStartEvent)
				c, err := di.getSingletone(event.alias)
				di.getContainerEndCh <- operationEndEvent{
					operationIndex,
					getContainerEndEvent{
						container: c,
						err:       err,
					},
				}
			case getTemporalOperation:
				event := startEvent.event.(getTemporalStartEvent)
				c, err := di.getTemporal(event.alias, event.args...)
				di.getContainerEndCh <- operationEndEvent{
					operationIndex,
					getContainerEndEvent{
						container: c,
						err:       err,
					},
				}
			case runOperation:
				event := startEvent.event.(runStartEvent)
				cr, err := di.run(event.function, event.args...)
				di.runEndCh <- operationEndEvent{
					operationIndex,
					runEndEvent{
						funcEndEvent{
							cr:  cr,
							err: err,
						},
					},
				}
			case recoverOperation:
				event := startEvent.event.(recoverStartEvent)
				cr, err := di.recover(event.function, event.args...)
				di.recoverEndCh <- operationEndEvent{
					operationIndex,
					recoverEndEvent{
						funcEndEvent{
							cr:  cr,
							err: err,
						},
					},
				}
			case recoverAndCleanOperation:
				event := startEvent.event.(recoverAndCleanStartEvent)
				cr, err := di.recoverAndClean(event.function, event.args...)
				di.recoverAndCleanEndCh <- operationEndEvent{
					operationIndex,
					recoverAndCleanEndEvent{
						funcEndEvent{
							cr:  cr,
							err: err,
						},
					},
				}
			}
		case <-di.exitCh:
			return
		}
	}
}
