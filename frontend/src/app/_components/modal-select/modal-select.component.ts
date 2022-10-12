import {Component, OnInit} from '@angular/core';
import {MdbModalRef} from "mdb-angular-ui-kit/modal";
import {Role} from "../../models/role.model";
import {Observable} from "rxjs";

@Component({
  selector: 'app-modal-select',
  templateUrl: './modal-select.component.html',
  styleUrls: ['./modal-select.component.scss']
})
export class ModalSelectComponent {
  ignoreOptions: any[];
  options: any[];
  selectedOptions: any[] = [];

  title: string | null;
  message: string | null;

  constructor(public modalRef: MdbModalRef<ModalSelectComponent>) {
  }

  closeWithRoles() {
    this.modalRef.close(this.selectedOptions);
  }
}
