import './Home.css';
import { Fragment, useEffect, useState } from 'react';
import { RiAccountCircleFill } from 'react-icons/ri';
import { Menu, Transition } from '@headlessui/react';
import { useNavigate } from 'react-router-dom';

function classNames(...classes: string[]) {
  return classes.filter(Boolean).join(' ');
}

export default function Home() {
  const [loading, setLoading] = useState(true);
  const [progress, setProgress] = useState(0.0);
  const [dlSpeed, setDlSpeed] = useState(0.0);
  const [dlUnits, setDLUnits] = useState('KB/s');
  const [downloading, setDownloading] = useState(false);
  const [isInstalled, setIsInstalled] = useState(true);

  const navigate = useNavigate();

  useEffect(() => {
    window.electron.ipcRenderer.sendMessage('game-status', []);

    window.electron.ipcRenderer.on('installed', () => {
      setIsInstalled(true);
    });

    window.electron.ipcRenderer.once('game-status', (installed) => {
      setIsInstalled(!!installed);
      setLoading(false);
    });

    window.electron.ipcRenderer.once('download', () => {
      setDownloading(false);
    });

    window.electron.ipcRenderer.on('download-progress', (p: unknown) => {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      setProgress(Math.floor(p.percent * 100));
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      setDlSpeed(Math.round(p.speed * 10) / 10);
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      setDLUnits(p.units);
    });
  }, []);

  const onSignOut = () => {
    navigate('/');
  };

  const onInstall = () => {
    setIsInstalled(false);
    setDownloading(true);
    setProgress(0.0);
    setDlSpeed(0.0);
    console.log('renderer call download');
    window.electron.ipcRenderer.sendMessage('download', []);
  };

  const onCancelDownload = () => {
    window.electron.ipcRenderer.sendMessage('download-cancel', []);
    setDownloading(false);
  };

  const onPlay = () => {
    window.electron.ipcRenderer.sendMessage('launch-client', []);
  };

  return (
    <div className="main-bg w-full h-screen text-center">
      <div className="logo absolute top-[100px] left-[450px] w-[570px] h-[150px]" />
      <Menu
        as="div"
        className="absolute top-5 left-5 rounded-full w-[30px] h-[30px]"
      >
        <div>
          <Menu.Button>
            <RiAccountCircleFill className="w-full h-full text-gray-200 text-3xl hover:text-blue-300" />
            {/* <div className="relative rounded-full bottom-3 left-0 w-[10px] h-[10px] bg-green-300 shadow shadow-black" /> */}
          </Menu.Button>
        </div>

        <Transition
          as={Fragment}
          enter="transition ease-out duration-100"
          enterFrom="transform opacity-0 scale-95"
          enterTo="transform opacity-100 scale-100"
          leave="transition ease-in duration-75"
          leaveFrom="transform opacity-100 scale-100"
          leaveTo="transform opacity-0 scale-95"
        >
          <Menu.Items className="absolute left-0 z-10 mt-2 w-56 origin-top-right bg-gray-800 rounded-md shadow-lg ring-1 ring-white ring-opacity-5 focus:outline-none">
            <div className="py-1">
              <Menu.Item>
                {({ active }) => (
                  <a
                    href="https://shatteredrealmsonline.com/"
                    rel="noreferrer"
                    target="_blank"
                    className={classNames(
                      active ? 'bg-gray-900 text-gray-100' : 'text-gray-400',
                      'block px-4 py-2 text-sm'
                    )}
                  >
                    Profile
                  </a>
                )}
              </Menu.Item>
              <Menu.Item>
                {({ active }) => (
                  // eslint-disable-next-line jsx-a11y/anchor-is-valid
                  <button
                    type="button"
                    className={classNames(
                      active ? 'bg-gray-900 text-gray-100' : 'text-gray-400',
                      'block w-full px-4 py-2 text-sm'
                    )}
                    onClick={onSignOut}
                  >
                    Sign out
                  </button>
                )}
              </Menu.Item>
            </div>
          </Menu.Items>
        </Transition>
      </Menu>
      <div className="absolute left-0 bottom-0 w-full h-[200px]">
        <div className="w-11/12 h-full inline-flex items-center justify-center m-0">
          <div className="flex w-full h-[100px] items-center">
            <div
              className="w-3/4 h-[10px] bg-blue-500/30 shadow-xl shadow-blue-300 rounded-l"
              style={{ opacity: downloading ? '100%' : '0' }}
            >
              <div
                className="h-full bg-blue-500 rounded-l"
                style={{ width: `${progress}%` }}
              />
            </div>
            <div className="flex w-1/4 h-full bg-gray-800">
              {loading ? (
                <svg
                  className="animate-spin p-5 h-full w-full text-white"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  />
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
              ) : (
                <>
                  {isInstalled ? (
                    <button
                      onClick={onPlay}
                      type="button"
                      className="w-full text-5xl text-blue-300 hover:bg-gray-900"
                    >
                      Play
                    </button>
                  ) : (
                    <>
                      {downloading ? (
                        <div className="flex w-full h-full items-center">
                          <div className="w-10/12 text-4xl text-center text-gray-300">
                            {progress}%
                            <p className="text-xs text-gray-500">
                              Downloading - {dlSpeed} {dlUnits}
                            </p>
                          </div>
                          <button
                            className="w-2/12 h-full text-center bg-red-900 hover:bg-red-800"
                            onClick={onCancelDownload}
                            type="button"
                          >
                            X
                          </button>
                        </div>
                      ) : (
                        <>
                          {progress < 100 ? (
                            <button
                              onClick={onInstall}
                              type="button"
                              className="w-full text-5xl text-blue-300 hover:bg-gray-900"
                            >
                              Install
                            </button>
                          ) : (
                            <button
                              className="w-full h-full text-5xl text-gray-300 cursor-default"
                              type="button"
                            >
                              Installing
                            </button>
                          )}
                        </>
                      )}
                    </>
                  )}
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
