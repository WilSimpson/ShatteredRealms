import {Component} from '@angular/core';
import { MdbNotificationRef } from 'mdb-angular-ui-kit/notification';

@Component({
  selector: 'app-alert',
  templateUrl: './alert.component.html',
  styleUrls: ['./alert.component.scss']
})
export class AlertComponent {
  message: string | null = null;
  closeable: boolean = true;
  color: 'primary' |
    'secondary' |
    'success' |
    'danger' |
    'warning' |
    'info' |
    'light' |
    'dark' = 'primary';
  fade: boolean = false;

  constructor(public notificationRef: MdbNotificationRef<AlertComponent>) {}

  getHTMLClasses(): string {
    let classes = `alert alert-${this.color} show`;

    if (this.closeable) classes +=' alert-dismissible'
    if (this.fade) classes +=' fade'

    return classes;
  }
}
