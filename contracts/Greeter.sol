pragma solidity ^0.4.18;

//
// From: https://github.com/ethereum/go-ethereum/wiki/Contract-Tutorial
//
// The Greeter is an intelligent digital entity that lives on the blockchain and is able to have conversations with
// anyone who interacts with it, based on its input. It might not be a talker, but it’s a great listener. Here is its
// code:
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

contract Greeter is mortal {
	string greeting;
	uint256 public ngreeting;
	string public xgreeting;

    event ReportGreetingEvent(string greeting);
    event ReportGreetingChangedEvent(string greeting);
	
	function greeter() public {
		greeting = "not set";
		xgreeting = "x not set";
	}

	// this runs when the contract is executed 
	function greeter(string _greeting) public {
		greeting = _greeting;
		xgreeting = "x is the X";
	}

	// Create a new greeting
	function setGreeting(string _greeting) public {
		greeting = _greeting;
		ReportGreetingChangedEvent(greeting);
	}

	function setNGreeting(uint256 _greeting) public {
		ngreeting = _greeting;
	}

	// Get the data back as an event.
	function getGreeting() public returns (string) {
    	ReportGreetingEvent(greeting);
		return greeting;
	}
	function getNGreeting() public view returns (uint256) {
		return ngreeting;
	}

	function test01(uint256 aNum) public payable returns ( uint256 id ) {
		if ( aNum > 100 ) {
			ngreeting = ngreeting + 1;
		} else {
			ngreeting = ngreeting - 1;
		}
		return ngreeting;
	}

	// main function 
	function greet() public constant returns (string) {
		return greeting;
	}
}
