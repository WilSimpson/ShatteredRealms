import axios from 'axios';

let ACCOUNTS_API = '';

window.electron.ipcRenderer.on('accounts-api-url', (resp) => {
  if (typeof resp === 'string') {
    ACCOUNTS_API = resp;
  } else {
    ACCOUNTS_API = 'https://api.shatteredrealmsonline.com/accounts';
  }
});

window.electron.ipcRenderer.sendMessage('accounts-api-url', []);

const AccountsService = {
  login: (email: string, password: string) => {
    return axios.post(`${ACCOUNTS_API}/v1/login`, { email, password });
  },
};

export default AccountsService;
