import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {Router} from "@angular/router";
import {AuthenticationService} from "../../_services/authentication.service";
import {environment} from "../../../environments/environment";
import {ReCaptchaV3Service} from "ngx-captcha";
import {AlertComponent} from "../../_components/alert/alert.component";
import {MdbNotificationService} from "mdb-angular-ui-kit/notification";

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss']
})
export class RegisterComponent implements OnInit {
  registerForm: FormGroup = new FormGroup({
    firstName: new FormControl('', [Validators.required, Validators.maxLength(50)]),
    lastName: new FormControl('', [Validators.required, Validators.maxLength(50)]),
    email: new FormControl('', [Validators.required, Validators.email, Validators.maxLength(100)]),
    username: new FormControl('', [Validators.required, Validators.minLength(3), Validators.max(25)]),
    password: new FormControl('', [Validators.required, Validators.minLength(6), Validators.maxLength(64)]),
    confirmPassword: new FormControl('', [Validators.required, Validators.minLength(6), Validators.maxLength(64)]),
  });
  loading: boolean;
  siteKey: string;

  constructor(protected router: Router,
              protected authService: AuthenticationService,
              protected notificationService: MdbNotificationService,
              protected reCaptchaV3Service: ReCaptchaV3Service) {
  }

  ngOnInit(): void {
    this.siteKey = environment.recaptcha.siteKey;
    if (this.authService.currentUserValue) {
      this.router.navigate(['/']).then(() => {
        this.notificationService.open(AlertComponent, {
          data: {
            message: 'You are already signed in.',
            color: 'warning',
            fade: true,
          },
          stacking: true,
          position: "top-center",
        })
      });

      return;
    }
  }

  onRegister(): void {
    if (!this.registerForm.valid) {
      this.notificationService.open(AlertComponent, {
        data: {
          message: 'Please fill out all forms correctly',
          color: 'warning',
          fade: true,
        },
        stacking: true,
        position: "top-center",
      })
      return;
    }

    if (this.registerForm.value.password !== this.registerForm.value.confirmPassword) {
      this.notificationService.open(AlertComponent, {
        data: {
          message: 'Passwords do not match',
          color: 'warning',
          fade: true,
        },
        stacking: true,
        position: "top-center",
      })
      return;
    }

    this.reCaptchaV3Service.execute(this.siteKey, 'homepage', (token) => {
      console.log('This is your token: ', token);
    }, {
      useGlobalDomain: false
    });

    this.authService.requestRegister({
      email: this.registerForm.value.email,
      password: this.registerForm.value.password,
      firstName: this.registerForm.value.firstName,
      lastName: this.registerForm.value.lastName,
      username: this.registerForm.value.username,
    }).subscribe((success) => {
      if (success.errors) {
        for (const e of success.errors) {
          this.notificationService.open(AlertComponent, {
            data: {
              message: `${e.text}: ${e.info}`,
              color: 'warning',
            },
            stacking: true,
            position: "top-center",
          })
        }
      } else {
        this.router.navigate(['/']).then(() => {
          this.notificationService.open(AlertComponent, {
            data: {
              message: 'Account created!',
              color: 'success',
            },
            stacking: true,
            position: "top-center",
          })
        });
      }
    }, (error) => {
      this.notificationService.open(AlertComponent, {
        data: {
          message: error.error.message,
          color: 'danger',
        },
        stacking: true,
        position: "top-center",
      })
    }).add(() => {
      this.loading = false;
    });
  }
}
