import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {MdbTableDirective} from "mdb-angular-ui-kit/table";
import {User} from "../../models/user.model";
import {UserPermission} from "../../models/user-permission.model";
import {MdbCheckboxChange} from "mdb-angular-ui-kit/checkbox";

@Component({
  selector: 'app-permissions-table',
  templateUrl: './permissions-table.component.html',
  styleUrls: ['./permissions-table.component.scss']
})
export class PermissionsTableComponent implements OnInit {
  @ViewChild('table') table!: MdbTableDirective<User>;

  @Input() dataSource: UserPermission[];
  @Input() loading = true;
  @Input() showActions = false;
  @Input() searchable = false;
  @Input() selectable = false;

  @Input() selected = new Set<UserPermission>();

  @Output() outSelected = new EventEmitter<Set<UserPermission>>();

  constructor() {
  }

  ngOnInit(): void {
  }

  updateOutput() {
    this.outSelected.emit(this.selected);
  }

  search(searchTerm: string): void {
    this.table.search(searchTerm);
  }

  areAllChecked() {
    const selectionsLength = this.selected.size;
    const dataLength = this.dataSource.length;
    return selectionsLength === dataLength;  }

  toggleSelection(event: MdbCheckboxChange, value: UserPermission): void {
    if (event.checked) {
      this.select(value);
    } else {
      this.deselect(value);
    }

    this.updateOutput();
  }

  toggleAll(event: MdbCheckboxChange) {
    if (event.checked) {
      this.dataSource?.forEach((p: UserPermission) => {
        this.select(p);
      })
    } else {
      this.dataSource?.forEach((p: UserPermission) => {
        this.deselect(p);
      })
    }

    this.updateOutput();
  }

  select(value: UserPermission): void {
    if (!this.selected.has(value)) {
      this.selected.add(value);
    }

    this.updateOutput();
  }

  deselect(value: UserPermission): void {
    if (this.selected.has(value)) {
      this.selected.delete(value);
    }

    this.updateOutput();
  }
}
