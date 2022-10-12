import {Component, OnInit} from '@angular/core';
import {getRoleBadgeClasses, Role} from "../../../models/role.model";
import {AuthorizationService} from "../../../_services/authorization.service";
import {UserPermission} from "../../../models/user-permission.model";
import {ActivatedRoute} from "@angular/router";
import {MdbNotificationService} from "mdb-angular-ui-kit/notification";
import {AlertComponent} from "../../../_components/alert/alert.component";

@Component({
  selector: 'app-edit-role',
  templateUrl: './edit-role.component.html',
  styleUrls: ['./edit-role.component.scss']
})
export class EditRoleComponent implements OnInit {

  role: Role;
  permissions: UserPermission[] = [];
  selected = new Set<UserPermission>();

  newName: string = '';

  constructor(
    protected notificationService: MdbNotificationService,
    private authorizationService: AuthorizationService,
    private route: ActivatedRoute,
  ) {
  }

  ngOnInit(): void {
    this.authorizationService.getRole(Number(this.route.snapshot.paramMap.get('role'))).subscribe(
      (role) => {
        this.role = role;
        this.newName = role.name;
        this.setupPermissions();
      }
    );
  }

  setSelected(selected: Set<UserPermission>) {
    this.selected = selected;
  }

  setupPermissions() {
    this.authorizationService.getAllPermissions().subscribe(permissions => {
      permissions.forEach(permission => {
        let found = this.role.permissions.find(rp => rp.permission == permission.permission);
        if (found) {
          this.permissions.push(found);
          this.selected.add(found);
        } else {
          this.permissions.push({permission: permission.permission, other: false});
        }
      });
      this.permissions = [...this.permissions];
    });
  }

  getRoleClasses(): string {
    if (this.role)
      return getRoleBadgeClasses(this.role);

    return '';
  }

  submitNameChange() {
    this.authorizationService.updateRoleName(this.role.id, this.newName).subscribe({
      next: () => {
        this.role.name = this.newName;
        this.notificationService.open(AlertComponent, {
          data: {
            message: 'Name updated',
            color: 'danger',
          },
          stacking: true,
          autohide: true,
        });
      },
      error: err => {
        this.notificationService.open(AlertComponent, {
          data: {
            message: `Error: ${err.message}`,
            color: 'danger',
          },
          stacking: true
        });
      }
    });
  }

  cancelNameChange() {
    this.newName = this.role.name;
  }

  saveChanges() {
    this.authorizationService.updateRolePermissions(this.role.id, Array.from(this.selected)).subscribe({
      next: () => {
        this.notificationService.open(AlertComponent, {
          data: {
            message: 'Permissions updated',
            color: 'danger',
            fade: true,
          },
          stacking: true,
        });
      },
      error: err => {
        this.notificationService.open(AlertComponent, {
          data: {
            message: `Error: ${err.message}`,
            color: 'danger',
          },
          stacking: true
        });
      }
    });
  }
}
