import {ComponentFixture, TestBed, waitForAsync} from '@angular/core/testing';

import { LoginComponent } from './login.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AlertService} from "../../../_services/alert.service";
import {AuthenticationService} from "../../_services/authentication.service";
import {Observable, of, throwError} from "rxjs";
import {RouterTestingModule} from "@angular/router/testing";
import {User} from "../../models/user.model";
import {Role} from "../../models/role.model";
import {Router} from "@angular/router";
import {NewUser} from "../../_helpers/dummy.data";

describe('LoginComponent', () => {
  let component: LoginComponent;
  let fixture: ComponentFixture<LoginComponent>;
  let alertService: AlertService;
  let authService: AuthenticationService;
  let router: Router;
  let user: User;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [
        LoginComponent
      ],
      imports: [
        HttpClientTestingModule,
        RouterTestingModule,
      ]
    })
    .compileComponents();
    user = NewUser();
  });

  it('should redirect if a user is logged in', () => {
    authService = TestBed.inject(AuthenticationService);
    router = TestBed.inject(Router);

    spyOnProperty(authService, 'currentUserValue').and.returnValue({});
    spyOn(router, 'navigate').and.returnValue(Promise.resolve(true));

    fixture = TestBed.createComponent(LoginComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();

    expect(router.navigate).toHaveBeenCalled();
  });

  describe('when a user is not signed in', () => {
    beforeEach(() => {
      fixture = TestBed.createComponent(LoginComponent);
      component = fixture.componentInstance;
      fixture.detectChanges();
      alertService = TestBed.inject(AlertService);
      authService = TestBed.inject(AuthenticationService);
      router = TestBed.inject(Router);
    });

    it('should create', () => {
      expect(component).toBeTruthy();
      expect(component.loginForm).toBeTruthy();
    });

    describe('login', () => {
      let validUser: any;

      beforeEach(() => {
        validUser = {
          email: 'email@email.com',
          password: 'password',
        };
        component.loginForm.setValue(validUser);
      })

      describe('valid input', () => {
        it('should be truthy', () => {
          expect(component.onLogin()).toBeTruthy();
        });

        it('should make a request to login', () => {
          spyOn(authService, 'login').and.returnValue(of(user));
          component.onLogin();
          expect(authService.login).toHaveBeenCalledOnceWith(validUser.email, validUser.password);
          expect(authService.login).toHaveBeenCalledTimes(1);
        });

        it('should login the user on success', () => {
          spyOn(authService, 'login').and.returnValue(of(user));
          const promise = Promise.resolve(true);
          spyOn(router, 'navigate').and.returnValue(promise);
          spyOn(alertService, 'success');

          component.onLogin();

          expect(authService.login).toHaveBeenCalledOnceWith(validUser.email, validUser.password);
          expect(authService.login).toHaveBeenCalledTimes(1);

          expect(router.navigate).toHaveBeenCalledOnceWith(['/']);
          expect(router.navigate).toHaveBeenCalledTimes(1);

          waitForAsync(() => expect(alertService.success).toHaveBeenCalledTimes(1));
        });

        it('should show an error on server failure', () => {
          const errorMessage = 'error';
          spyOn(authService, 'login').and.returnValue(throwError({message: errorMessage}));
          spyOn(router, 'navigate');
          spyOn(alertService, 'error');

          component.onLogin();

          expect(authService.login).toHaveBeenCalledOnceWith(validUser.email, validUser.password);
          expect(authService.login).toHaveBeenCalledTimes(1);

          expect(router.navigate).toHaveBeenCalledTimes(0);
          expect(alertService.error).toHaveBeenCalledTimes(1);
        });

        it('should show an error on invalid credentials', () => {
          spyOn(authService, 'login').and.returnValue(throwError({status: 401}));
          spyOn(router, 'navigate');
          spyOn(alertService, 'warn');

          component.onLogin();

          expect(authService.login).toHaveBeenCalledOnceWith(validUser.email, validUser.password);
          expect(authService.login).toHaveBeenCalledTimes(1);

          expect(router.navigate).toHaveBeenCalledTimes(0);
          expect(alertService.warn).toHaveBeenCalledTimes(1);
        });
      });


      it('should validate the email', () => {
        validUser.email = "asdf";
        component.loginForm.setValue(validUser);
        expect(component.onLogin()).toBeFalsy();
      });

      it('should validate the password', () => {
        validUser.password = "";
        component.loginForm.setValue(validUser);
        expect(component.onLogin()).toBeFalsy();
      });
    });
  });
});
