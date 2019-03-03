require('babel-register');
require('babel-polyfill');

module.exports = {
	networks: {
		development: {
			host: "127.0.0.1",
			port: 9545,
			network_id: "*", // Match any network id
			gas: 6000000,
			from: "0x627306090abab3a6e1400e9345bc60c78a8bef57"
		},
		local: {
			host: "192.168.0.158",
			port: 8545,
			network_id: "91204",
			gas: 4712388,
			from: "0xc2a56884538778bacd91aa5bf343bf882c5fb18b"
		}
	}
};
