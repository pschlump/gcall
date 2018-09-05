pragma solidity ^0.4.18;

//
// From: https://github.com/ethereum/go-ethereum/wiki/Contract-Tutorial
//

//$version$: v1.0.0
//$author$: Philip Schlump
//$date$: Tue Mar 13 16:25:46 MDT 2018

contract mortal {
	/* Define variable owner of the type address*/
	address public owner = msg.sender;

	/* this function is executed at initialization and sets the owner of the contract */
	function mortal() public {
		// owner = msg.sender;
	}

	/* Function to recover the funds on the contract */
	function kill() public {
		if (msg.sender == owner) {
			// suicide(owner);
			selfdestruct(owner);
		}
	}

	// To be able to move payments from this contract back to Keep multiwallet
	function withdraw() public {
		require(owner == msg.sender);
		owner.transfer(this.balance);
	}

	// To be able to move payments from this contract back to Keep multiwallet
	function withdrawAmount(uint256 _amt) public {
		require(owner == msg.sender);
		if ( _amt <= this.balance ) {
			owner.transfer(_amt);
		}
	}
}

