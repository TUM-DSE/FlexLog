START_SEQUENCER = ./cmd/start_sequencer/start_sequencer.go
SEQUENCER_BIN = ./seq

seq:
	go build -o $(SEQUENCER_BIN) $(START_SEQUENCER)

clean:
	rm -f $(SEQUENCER_BIN)
	if pgrep $(SEQUENCER_BIN); then pkill $(SEQUENCER_BIN); fi