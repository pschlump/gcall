pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "../contracts/GroupProc.sol";

contract TestGroupProc {	
	
	GroupProc gp = new GroupProc();

	function testCreateAGroup() public {
		bytes32 _groupPubKey ;
		uint n;

		// After reset - should have 0 groups
		gp.ResetNoOfGroups();
		n = gp.GettNoOfGroups();
		Assert.equal(n,0,"should have 0 groups, did not");

		_groupPubKey = hex"0101";
		gp.GroupCreated( _groupPubKey ) ;
		n = gp.GettNoOfGroups();
		Assert.equal(n,1,"should have 1 group, did not");

		_groupPubKey = hex"0102";
		gp.GroupCreated( _groupPubKey ) ;
		n = gp.GettNoOfGroups();
		Assert.equal(n,2,"should have 2 groups, did not");

		// Remove non-existent group - should have same number of groups
		_groupPubKey = hex"0104";
		gp.GroupDistroyed( _groupPubKey ) ;
		n = gp.GettNoOfGroups();
		Assert.equal(n,2,"should have 2 groups, did not");

		// Add group twice - should have non-existent group
		_groupPubKey = hex"0102";
		gp.GroupCreated( _groupPubKey ) ;
		n = gp.GettNoOfGroups();
		Assert.equal(n,2,"should have 2 groups, did not");

		_groupPubKey = hex"0101";
		gp.GroupDistroyed( _groupPubKey ) ;
		n = gp.GettNoOfGroups();
		Assert.equal(n,1,"should have 1 group, did not");

		// After reset - should have 0 groups
		gp.ResetNoOfGroups();
		n = gp.GettNoOfGroups();
		Assert.equal(n,0,"should have 0 groups, did not");



		// TODO - in .js test (catch events)
	}

}
