## Таблица с математическими формулами

| Формула | Описание | Применение |
|---------|----------|------------|
| $E = mc^2$ | Эквивалентность массы и энергии | Ядерная физика |
| $F = G\frac{m_1m_2}{r^2}$ | Закон всемирного тяготения | Астрономия |
| $\vec{E} = \rho / \varepsilon_0$ | Уравнение Максвелла | Электродинамика |
| $\Delta S \geq 0$ | Второй закон термодинамики | Термодинамика |

## Блоки кода с комментариями о формулах

```python
# Вычисление квадратного корня (√x)
import math
x = 16
sqrt_x = math.sqrt(x)  # Эквивалентно x^(1/2)
print(f"Квадратный корень из {x} равен {sqrt_x}")

# Вычисление факториала (n!)
def factorial(n):
    """
    Вычисляет n! (произведение натуральных чисел от 1 до n)
    
    В математике: n! = ∏_{i=1}^{n} i
    """
    if n == 0:
        return 1
    result = 1
    for i in range(1, n+1):
        result *= i
    return result

print(f"5! = {factorial(5)}")
```
