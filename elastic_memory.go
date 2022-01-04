package gmachine

type ElasticMemory map[Word]Word

func (e ElasticMemory) Add(words []Word) {
	length := len(e)

	memoryLocation := 0

	if length > 0 {
		memoryLocation += 1
	}

	for index, word := range words {
		e[Word(memoryLocation+index)] = word
	}

}

func (e ElasticMemory) AddToAddress(address Word, word Word) {

	e[address] = word

}
