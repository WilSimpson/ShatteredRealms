import {Component} from '@angular/core';
import {MdbModalRef} from 'mdb-angular-ui-kit/modal';

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
})
export class ModalComponent {
  title: string | null = null;
  message: string | null = null;
  closeText: string = 'Cancel';
  submitText: string = 'Save Changes';

  constructor(public modalRef: MdbModalRef<ModalComponent>) {
  }
}
