package common

func vote() {}

func sendTransList() {

}

func sendBlock() {}

func sendBlockRound2() {}

func voteForBlockRound1() {}

func voteForBlockRound2() {}

func MaxByzantiumNumber(total int) int {
	if total % 3 == 0 {
		return total / 3 - 1
	} else {
		return total / 3
	}
}

func QuorumNumber(total int) int {
	return 2 * MaxByzantiumNumber(total) + 1
}