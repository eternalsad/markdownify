# Сложные математические конструкции

## Матрицы и определители

Матрица $A$ размера $m \times n$:

$$A = \begin{pmatrix}
a_{11} & a_{12} & \ldots & a_{1n} \\
a_{21} & a_{22} & \ldots & a_{2n} \\
\vdots & \vdots & \ddots & \vdots \\
a_{m1} & a_{m2} & \ldots & a_{mn}
\end{pmatrix}$$

Определитель матрицы $3 \times 3$:

$$\det(A) = \begin{vmatrix}
a_{11} & a_{12} & a_{13} \\
a_{21} & a_{22} & a_{23} \\
a_{31} & a_{32} & a_{33}
\end{vmatrix} = a_{11}(a_{22}a_{33} - a_{23}a_{32}) - a_{12}(a_{21}a_{33} - a_{23}a_{31}) + a_{13}(a_{21}a_{32} - a_{22}a_{31})$$

## Системы уравнений

Система линейных уравнений:

$$\begin{cases}
a_{11}x_1 + a_{12}x_2 + \ldots + a_{1n}x_n = b_1 \\
a_{21}x_1 + a_{22}x_2 + \ldots + a_{2n}x_n = b_2 \\
\vdots \\
a_{m1}x_1 + a_{m2}x_2 + \ldots + a_{mn}x_n = b_m
\end{cases}$$

## Комплексные числа

Комплексное число $z$ в алгебраической форме:

$$z = a + bi$$

где $i^2 = -1$

Формула Эйлера:

$$e^{ix} = \cos x + i\sin x$$

## Ряды и суммы

Ряд Тейлора для функции $f(x)$ в окрестности точки $a$:

$$f(x) = f(a) + \frac{f'(a)}{1!}(x-a) + \frac{f''(a)}{2!}(x-a)^2 + \frac{f'''(a)}{3!}(x-a)^3 + \ldots$$

Сумма арифметической прогрессии:

$$S_n = \frac{n(a_1 + a_n)}{2}$$

Сумма геометрической прогрессии:

$$S_n = \frac{a_1(1-q^n)}{1-q}, q \neq 1$$

## Дифференциальные уравнения

Линейное дифференциальное уравнение второго порядка:

$$a\frac{d^2y}{dx^2} + b\frac{dy}{dx} + cy = f(x)$$

Волновое уравнение:

$$\frac{\partial^2 u}{\partial t^2} = c^2 \frac{\partial^2 u}{\partial x^2}$$

Уравнение теплопроводности:

$$\frac{\partial u}{\partial t} = \alpha \frac{\partial^2 u}{\partial x^2}$$

## Многомерные интегралы

Двойной интеграл:

$$\iint_D f(x,y) \, dA = \iint_D f(x,y) \, dx \, dy$$

Тройной интеграл:

$$\iiint_V f(x,y,z) \, dV = \iiint_V f(x,y,z) \, dx \, dy \, dz$$