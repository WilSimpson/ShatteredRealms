import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {MdbTableDirective} from "mdb-angular-ui-kit/table";
import {getRoleBadgeClasses, Role} from "../../models/role.model";
import {Router} from "@angular/router";
import {AuthorizationService} from "../../_services/authorization.service";
import {MdbModalService} from "mdb-angular-ui-kit/modal";
import {ModalComponent} from "../modal/modal.component";
import {MdbNotificationService} from "mdb-angular-ui-kit/notification";
import {AlertComponent} from "../alert/alert.component";

@Component({
  selector: 'app-roles-table',
  templateUrl: './roles-table.component.html',
  styleUrls: ['./roles-table.component.scss']
})
export class RolesTableComponent implements OnInit {
  @ViewChild('table') table!: MdbTableDirective<Role>;

  @Input() dataSource: Role[] | null;
  @Input() loading = true;
  @Input() showActions = false;
  @Input() searchable = false;
  @Input() rowSelectable = true;

  constructor(
    private router: Router,
    private authorizationService: AuthorizationService,
    private modalService: MdbModalService,
    private notificationService: MdbNotificationService,
  ) {
  }

  ngOnInit(): void {
  }

  search(searchTerm: string): void {
    this.table.search(searchTerm);
  }

  getRoleHTMLClass(role: Role): string {
    return getRoleBadgeClasses(role);
  }

  roleClicked(role: Role) {
    if (this.rowSelectable) {
      this.router.navigate(['/roles', role.id])
    }
  }

  roleEditClicked(role: Role) {
    this.router.navigate(['/roles', role.id, 'edit'])
  }

  roleDeleteClicked(role: Role) {
    this.modalService.open(ModalComponent, {
      data: {
        title: 'Deletion confirmation',
        message: `Are you sure you want to delete the role '${role.name}'`,
        submitText: 'Delete',
      }
    }).onClose.subscribe((message) => {
      if (message) {
        this.authorizationService.deleteRole(role.id).subscribe({
          next: () => {
            this.notificationService.open(AlertComponent, {
              data: {
                message: 'Successfully deleted role',
                color: 'info',
              }
            });
            window.location.reload();
          },
          error: err => {
            this.notificationService.open(AlertComponent, {
              data: {
                message: `Error: ${err.error}`,
                color: 'danger',
              }
            });
          }
        });
      }
    });
  }
}
