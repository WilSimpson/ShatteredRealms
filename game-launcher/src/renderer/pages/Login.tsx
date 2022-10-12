import React, { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import AccountsService from '../../services/accounts.service';
import './Login.css';

export default function Login() {
  const [email, setEmail] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [errorMessage, setErrorMessage] = React.useState('');
  const navigate = useNavigate();

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setErrorMessage('');
    AccountsService.login(email, password)
      .then(() => {
        // return window.electron.ipcRenderer.sendMessage('navigate', ['/home']);
        // eslint-disable-next-line react/prop-types
        return navigate('/home');
      })
      .catch((error) => {
        if (error.response?.data) {
          setErrorMessage(error.response.data.message);
        } else {
          setErrorMessage('Unknown error. Please try again later.');
        }
      });
  }

  return (
    <div className="flex h-screen ml-auto w-[440px] bg-gray-800 shadow-2xl shadow-black">
      <div className="pt-[200px] w-80 h-full mx-auto text-center">
        <div className="logo h-[85px] mb-8" />
        {errorMessage ? (
          <div className="alert alert-danger">{errorMessage}</div>
        ) : null}
        <form onSubmit={handleSubmit}>
          <div className="form-input">
            <input
              type="email"
              placeholder="Email"
              name="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              required
            />
          </div>
          <div className="form-input">
            <a
              href="https://shatteredrealmsonline.com/"
              target="_blank"
              rel="noreferrer"
              className="float-right text-xs text-gray-400 hover:text-gray-300"
            >
              Forgot password
            </a>
            <input
              type="password"
              placeholder="Password"
              name="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              required
            />
          </div>
          <div className="form-input">
            <button
              type="submit"
              className="w-full bg-blue-500 p-2 my-2 text-center border border-gray-900 shadow-md shadow-gray-900 hover:bg-blue-400"
            >
              Login
            </button>
          </div>
        </form>
        <a
          href="https://shatteredrealmsonline.com/register"
          target="_blank"
          rel="noreferrer"
          className="w-full text-xs text-gray-500 hover:text-gray-300"
        >
          Create New Account
        </a>
      </div>
    </div>
  );
}
