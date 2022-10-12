import {Component, OnInit} from '@angular/core';
import {getRoleBadgeClasses, Role} from "../../../models/role.model";
import {AuthorizationService} from "../../../_services/authorization.service";
import {ActivatedRoute} from "@angular/router";
import {User} from "../../../models/user.model";
import {UsersService} from "../../../_services/users.service";

@Component({
  selector: 'app-view-role',
  templateUrl: './view-role.component.html',
  styleUrls: ['./view-role.component.scss']
})
export class ViewRoleComponent implements OnInit {

  role: Role;
  users: User[] = [];

  loadingUsers = true;

  constructor(
    private authorizationService: AuthorizationService,
    private usersService: UsersService,
    private route: ActivatedRoute,
  ) { }

  ngOnInit(): void {
    this.authorizationService.getRole(Number(this.route.snapshot.paramMap.get('role'))).subscribe(
      (role) => {
        this.role = role;
      }
    );

    this.usersService.getAllUsers().subscribe(
      users => {
        this.users = users.filter(u => u.roles.some(r => r.id == this.role.id));
        this.loadingUsers = false;
      }
    )
  }

  getRoleClasses(): string {
    if (this.role)
      return getRoleBadgeClasses(this.role);

    return '';
  }
}
