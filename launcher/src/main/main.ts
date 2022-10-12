/* eslint global-require: off, no-console: off, promise/always-return: off */

/**
 * This module executes inside of electron's main process. You can start
 * electron renderer process from here and communicate with the other processes
 * through IPC.
 *
 * When running `npm run build` or `npm run build:main`, this file is compiled to
 * `./src/main.js` using webpack. This gives us some performance wins.
 */
import path from 'path';
import { app, BrowserWindow, shell, ipcMain } from 'electron';
import { autoUpdater } from 'electron-updater';
import log from 'electron-log';
import { Progress } from 'electron-dl';
import { resolveHtmlPath } from './util';
import DownloadItem = Electron.DownloadItem;

const { download } = require('electron-dl');

const WindowsClientURL =
  'https://downloads.shatteredrealmsonline.com/client/WindowsClient.zip';

const AdmZip = require('adm-zip');
const fs = require('fs');

const child = require('child_process').execFile;

class AppUpdater {
  constructor() {
    log.transports.file.level = 'info';
    autoUpdater.logger = log;
    autoUpdater.checkForUpdatesAndNotify();
  }
}

let mainWindow: BrowserWindow | null = null;

let currentDownload: DownloadItem | null = null;
let lastTime = 0;
let lastSize = 0;

const installDirectory =
  process.env.NODE_ENV === 'production'
    ? `${path.dirname(app.getPath('exe'))}/game`
    : `${app.getAppPath()}/game`;

// ipcMain.on('ipc-example', async (event, arg) => {
//   const msgTemplate = (pingPong: string) => `IPC test: ${pingPong}`;
//   console.log(msgTemplate(arg));
//   event.reply('ipc-example', msgTemplate('pong'));
// });

ipcMain.on('minimize-window', async () => {
  mainWindow?.minimize();
});

ipcMain.on('navigate', async (_event, arg) => {
  mainWindow?.loadURL(resolveHtmlPath(`${arg}`));
});

ipcMain.on('accounts-api-url', async (event) => {
  event.reply('accounts-api-url', process.env.ACCOUNTS_API);
});

ipcMain.on('game-status', async (event) => {
  fs.access(`${installDirectory}/SRO`, (err: never) => {
    if (err) {
      // Game isn't installed
      event.reply('game-status', false);
    } else {
      event.reply('game-status', true);
    }
  });
});

ipcMain.on('launch-client', async () => {
  child(`${installDirectory}/SROClient.exe`, (err: never) => {
    if (!err) {
      app.quit();
      return;
    }

    console.log(err);
  });
});

ipcMain.on('download', async (event) => {
  download(mainWindow, WindowsClientURL, {
    directory: installDirectory,
    onStarted: (item: DownloadItem) => {
      currentDownload = item;
    },
    onProgress: (progress: Progress) => {
      const bytes = progress.transferredBytes - lastSize;
      lastSize = progress.transferredBytes;

      const ms = new Date().getTime() - lastTime;
      lastTime = new Date().getTime();

      let speed = bytes / ms;
      let units = 'KB/s';

      if (speed > 1000) {
        speed /= 1000;
        units = 'MB/s';
      }

      if (Number.isNaN(speed)) {
        speed = 0;
      }

      event.reply('download-progress', { ...progress, speed, units });
    },
    onCompleted: (file: File) => {
      event.reply('download', 'download completed');
      const zip = new AdmZip(file.path);
      zip.extractAllToAsync(installDirectory, true, true, () => {
        event.reply('installed', 'installation complete');
      });
      if (fs.existsSync(file.path)) {
        fs.unlink(file.path, () => {});
      }
    },
    onCancel: (item: DownloadItem) => {
      console.log('canceled: ', JSON.stringify(item));
    },
  });
});

ipcMain.on('download-cancel', async () => {
  if (currentDownload) {
    currentDownload.cancel();
    if (fs.existsSync(currentDownload.getSavePath())) {
      fs.unlink(currentDownload.getSavePath(), () => {});
    }
  }
});

if (process.env.NODE_ENV === 'production') {
  const sourceMapSupport = require('source-map-support');
  sourceMapSupport.install();
}

const isDebug =
  process.env.NODE_ENV === 'development' || process.env.DEBUG_PROD === 'true';

if (isDebug) {
  require('electron-debug')();
}

const installExtensions = async () => {
  const installer = require('electron-devtools-installer');
  const forceDownload = !!process.env.UPGRADE_EXTENSIONS;
  const extensions = ['REACT_DEVELOPER_TOOLS'];

  return installer
    .default(
      extensions.map((name) => installer[name]),
      forceDownload
    )
    .catch(console.log);
};

const createWindow = async () => {
  if (isDebug) {
    await installExtensions();
  }

  const RESOURCES_PATH = app.isPackaged
    ? path.join(process.resourcesPath, 'assets')
    : path.join(__dirname, '../../assets');

  const getAssetPath = (...paths: string[]): string => {
    return path.join(RESOURCES_PATH, ...paths);
  };

  mainWindow = new BrowserWindow({
    show: false,
    width: 1280,
    height: 720,
    icon: getAssetPath('icon.png'),
    webPreferences: {
      preload: app.isPackaged
        ? path.join(__dirname, 'preload.js')
        : path.join(__dirname, '../../.erb/dll/preload.js'),
    },
    fullscreenable: false,
    resizable: false,
    frame: false,
  });

  mainWindow.setTitle('Shattered Realms Online Launcher');

  mainWindow.loadURL(resolveHtmlPath('index.html'));

  mainWindow.on('ready-to-show', () => {
    if (!mainWindow) {
      throw new Error('"mainWindow" is not defined');
    }
    if (process.env.START_MINIMIZED) {
      mainWindow.minimize();
    } else {
      mainWindow.show();
    }
  });

  mainWindow.on('closed', () => {
    mainWindow = null;
  });

  // Open urls in the user's browser
  mainWindow.webContents.setWindowOpenHandler((edata) => {
    shell.openExternal(edata.url);
    return { action: 'deny' };
  });

  // Remove this if your app does not use auto updates
  // eslint-disable-next-line
  new AppUpdater();
};

/**
 * Add event listeners...
 */

app.on('window-all-closed', () => {
  // Respect the OSX convention of having the application in memory even
  // after all windows have been closed
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app
  .whenReady()
  .then(() => {
    createWindow();
    app.on('activate', () => {
      // On macOS it's common to re-create a window in the app when the
      // dock icon is clicked and there are no other windows open.
      if (mainWindow === null) createWindow();
    });
  })
  .catch(console.log);
