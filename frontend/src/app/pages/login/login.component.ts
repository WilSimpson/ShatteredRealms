import { Component, OnInit } from '@angular/core';
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {AuthenticationService} from "../../_services/authentication.service";
import {Router} from "@angular/router";
import {MdbNotificationService} from "mdb-angular-ui-kit/notification";
import {AlertComponent} from "../../_components/alert/alert.component";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  loginForm: FormGroup;
  loading: boolean = false;

  constructor(protected authService: AuthenticationService,
              public notificationService: MdbNotificationService,
              public router: Router) { }

  ngOnInit(): void {
    if (this.authService.currentUserValue) {
      this.router.navigate(['/']).then(() => {
        this.notificationService.open(AlertComponent, {
          data: {
            message: 'You are already signed in.',
            color: 'warning',
          },
          stacking: true,
          position: "top-center",
        })
      });
    }

    this.loginForm = new FormGroup({
      email: new FormControl('', [Validators.required, Validators.email]),
      password: new FormControl('', Validators.required),
    });
  }

  /**
   * Attempt to sign in the user with the given credentials.
   *
   * @return false if the login form has invalid data, otherwise true
   */
  onLogin(): boolean {
    if (!this.loginForm.valid) {
      this.notificationService.open(AlertComponent, {
        data: {
          message: 'Please fill out all fields correctly',
          color: 'warning',
        },
        stacking: true,
        position: "top-center",
      })
      return false;
    }
    this.loading = true;

    this.authService.login(this.loginForm.value.email, this.loginForm.value.password).subscribe((success) => {
      this.router.navigate(['/']).then(() => {
        this.notificationService.open(AlertComponent, {
          data: {
            message: 'Successfully logged in',
            color: 'success',
            fade: 'true',
          },
          stacking: true,
          position: "top-center",
        })
      });
    }, (error) => {
      if(error.status == 401) {
        this.notificationService.open(AlertComponent, {
          data: {
            message: 'Invalid email or password',
            color: 'warning',
          },
          stacking: true,
          position: "top-center",
        })
      } else {
        this.notificationService.open(AlertComponent, {
          data: {
            message: 'Error communicating with the server. Please try again later.',
            color: 'danger',
          },
          stacking: true,
          position: "top-center",
        })
      }
    }).add(() => {
      this.loading = false;
    });

    return true;
  }
}
