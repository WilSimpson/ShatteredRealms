import {Component, OnInit} from '@angular/core';
import {User} from "../../models/user.model";
import {Permission} from "../../models/permission.model";
import {ActivatedRoute, Router} from "@angular/router";
import {AuthenticationService} from "../../_services/authentication.service";
import {UsersService} from "../../_services/users.service";
import {AlertService} from "../../_services/alert.service";
import {Role} from "../../models/role.model";

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss']
})
export class UsersComponent implements OnInit {
  user: User = {
    created_at: new Date(Date.now()),
    email: "a",
    first_name: "",
    id: 0,
    last_name: "",
    permissions: [Permission.TEST_PERMISSION],
    role: Role.USER,
    token: "",
    username: ""
  };
  loading: boolean = true;

  constructor(protected route: ActivatedRoute,
              protected router: Router,
              protected alertService: AlertService,
              protected authService: AuthenticationService,
              protected usersService: UsersService) {}

  ngOnInit(): void {
    const id = Number(this.route.snapshot.paramMap.get('user'));
    if (this.authService.currentUserValue != null && this.authService.currentUserValue.id == id) {
      this.user = this.authService.currentUserValue;
    } else {
      this.usersService.getUser(id).subscribe((success) => {
        this.user = success.data;
      }, (error) => {
        this.router.navigate(['/']).then(() => {
          if (error.status == 404) {
            this.alertService.error('User not found');
          } else {
            this.alertService.error('Unknown server error. Please try again later.')
          }
        })
      })
    }
  }
}
