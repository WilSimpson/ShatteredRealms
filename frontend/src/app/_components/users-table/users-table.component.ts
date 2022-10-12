import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {User} from "../../models/user.model";
import {MdbTableDirective} from "mdb-angular-ui-kit/table";
import {getRoleBadgeClasses, Role} from "../../models/role.model";
import {advancedFilterFn} from "../../_helpers/filter.table";
import {UsersService} from "../../_services/users.service";
import {AlertComponent} from "../alert/alert.component";
import {MdbNotificationService} from "mdb-angular-ui-kit/notification";
import {Router} from "@angular/router";

@Component({
  selector: 'app-users-table',
  templateUrl: './users-table.component.html',
  styleUrls: ['./users-table.component.scss']
})
export class UsersTableComponent implements OnInit {
  @ViewChild('table') table!: MdbTableDirective<User>;

  @Input() dataSource: User[];
  @Input() loading = true;
  @Input() userSelectable = true;

  constructor(
    private usersService: UsersService,
    private notificationService: MdbNotificationService,
    private router: Router,
  ) {
  }

  ngOnInit(): void {
  }

  advancedSearch(value: string): void {
    this.table.search(value);
  }

  filterFn(data: any, searchTerm: string): boolean {
    return advancedFilterFn(data, searchTerm);
  }

  getRoleHTMLClasses(role: Role): string {
    return getRoleBadgeClasses(role);
  }

  banUser(id: any) {
    this.usersService.banUser(id).subscribe({
      next: () => {
        this.notificationService.open(AlertComponent, {
          data: {
            message: "User banned",
            color: 'info',
          },
          stacking: true,
          position: "top-center",
        });
        let user = this.dataSource.find(u => u.id == id);
        if (user) {
          user.bannedAt = new Date();
          this.dataSource = [...this.dataSource];
        }
      },
      error: (err) => {
        this.notificationService.open(AlertComponent, {
          data: {
            message: `ERROR: ${err.error.message}`,
            color: 'danger',
          },
          stacking: true,
          position: "top-center",
        });
      }
    })
  }

  unbanUser(id: any) {
    this.usersService.unBanUser(id).subscribe({
      next: () => {
        this.notificationService.open(AlertComponent, {
          data: {
            message: "User unbanned",
            color: 'success',
          },
          stacking: true,
          position: "top-center",
        });
        let user = this.dataSource.find(u => u.id == id);
        if (user) {
          user.bannedAt = null;
          this.dataSource = [...this.dataSource];
        }
      },
      error: (err) => {
        this.notificationService.open(AlertComponent, {
          data: {
            message: `ERROR: ${err.error.message}`,
            color: 'danger',
          },
          stacking: true,
          position: "top-center",
        });
      }
    })
  }

  isNaN(input: any): boolean {
   return isNaN(input);
  }

  onUserClick(user: User) {
    if (this.userSelectable) {
      this.router.navigate(['/users', user.id]);
    }
  }

  onUserEditClick(user: User) {
    this.router.navigate(['/users', user.id, 'edit']);
  }
}
