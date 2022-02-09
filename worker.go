package dilema

func (di *dicon) goQueueWriter() {
	for event := range di.queueCh {
		di.pushEventBack(event)
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
			operationCh := startEvent.operationCh
			var result operationEndEvent
			switch startEvent.oType {
			case registerTemporalOperation:
				event := startEvent.event.(registerTemporalStartEvent)
				err := di.registerTemporal(event.alias, event.serviceInit)
				result = operationEndEvent{
					registerEndEvent{err: err},
				}
			case registerSingleToneOperation:
				event := startEvent.event.(registerSingleToneStartEvent)
				err := di.registerSingleTone(event.alias, event.serviceInit, event.args...)
				result = operationEndEvent{
					registerEndEvent{err: err},
				}
			case registerFewOperation:
				event := startEvent.event.(registerFewStartEvent)
				err := di.registerFew(event.servicesInit, event.args...)
				result = operationEndEvent{
					registerEndEvent{err: err},
				}
			case getSingleToneOperation:
				event := startEvent.event.(getSingleToneStartEvent)
				c, err := di.getSingletone(event.alias)
				result = operationEndEvent{
					getContainerEndEvent{
						container: c,
						err:       err,
					},
				}
			case getTemporalOperation:
				event := startEvent.event.(getTemporalStartEvent)
				c, err := di.getTemporal(event.alias, event.args...)
				result = operationEndEvent{
					getContainerEndEvent{
						container: c,
						err:       err,
					},
				}
			case runOperation:
				event := startEvent.event.(runStartEvent)
				cr, err := di.run(event.function, event.args...)
				result = operationEndEvent{
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
				result = operationEndEvent{
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
				result = operationEndEvent{
					recoverAndCleanEndEvent{
						funcEndEvent{
							cr:  cr,
							err: err,
						},
					},
				}
			}

			operationCh <- result
		case <-di.exitCh:
			return
		}
	}
}
