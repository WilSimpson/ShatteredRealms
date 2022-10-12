import React from 'react';
import './default.scss';

interface TemplateProps {
  children: React.ReactNode;
}

const DefaultTemplate = (props: TemplateProps) => {
  function closeApp() {
    window.close();
  }

  function minimizeApp() {
    window.electron.ipcRenderer.sendMessage('minimize-window', []);
  }

  const { children } = props;

  return (
    <div>
      <header className="w-full absolute top-0 right-0 z-10 flex items-center justify-end">
        <button
          onClick={minimizeApp}
          className="hover:bg-white/10 w-[30px] h-[30px] pl-[10px]"
          type="button"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="11"
            height="2"
            viewBox="0 0 11 2"
          >
            <rect
              id="Rectangle_2"
              data-name="Rectangle 2"
              width="11"
              height="2"
              fill="#e2f7ff"
            />
          </svg>
        </button>
        <button
          onClick={closeApp}
          className="hover:bg-red-700 w-[30px] h-[30px] pl-[10px]"
          type="button"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="9.742"
            height="9.983"
            viewBox="0 0 9.742 9.983"
          >
            <path
              id="Icon_ionic-md-close"
              data-name="Icon ionic-md-close"
              d="M17.265,8.522l-.974-1-3.9,3.993L8.5,7.523l-.974,1,3.9,3.993-3.9,3.993.974,1,3.9-3.993,3.9,3.993.974-1-3.9-3.993Z"
              transform="translate(-7.523 -7.523)"
              fill="#e2f7ff"
            />
          </svg>
        </button>
      </header>
      {children}
    </div>
  );
};

export default DefaultTemplate;
