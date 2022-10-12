import { Component, OnInit } from '@angular/core';
import {UserPermission} from "../../../models/user-permission.model";
import {getRoleBadgeClasses, Role} from "../../../models/role.model";
import {AuthorizationService} from "../../../_services/authorization.service";
import {MdbNotificationService} from "mdb-angular-ui-kit/notification";
import {AlertComponent} from "../../../_components/alert/alert.component";
import {Router} from "@angular/router";

@Component({
  selector: 'app-new-role',
  templateUrl: './new-role.component.html',
  styleUrls: ['./new-role.component.scss']
})
export class NewRoleComponent implements OnInit {

  name: string = '';
  permissions: UserPermission[] = [];
  selected = new Set<UserPermission>();
  loading: boolean = true;

  constructor(
    private authorizationService: AuthorizationService,
    private notificationService: MdbNotificationService,
    private router: Router,
  ) { }

  ngOnInit(): void {
    this.authorizationService.getAllPermissions().subscribe(permissions => {
      permissions.forEach(permission => {
        this.permissions.push({permission: permission.permission, other: false});
      });
      this.permissions = [...this.permissions];
      this.loading = false;
    });
  }

  setSelected(selected: Set<UserPermission>) {
    this.selected = selected;
  }

  createRole() {
    if (this.name.length == 0) {
      this.notificationService.open(AlertComponent, {
        data: {
          message: 'Name cannot be empty',
          color: 'warning',
          fade: true,
        }
      });
      return
    }

    this.authorizationService.createRole(this.name, Array.from(this.selected)).subscribe({
      next: () => {
        this.router.navigate(['/admin']).then(() => {
          this.notificationService.open(AlertComponent, {
            data: {
              message: 'Role created!',
              color: 'success',
              fade: true,
            }
          });
        });
      },
      error: (err) => {
        console.log('error', err)
        this.notificationService.open(AlertComponent, {
          data: {
            message: `Error: ${err.message}`,
            color: 'danger',
            fade: true,
          }
        });
      }
    });
  }
}
