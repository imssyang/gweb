from typing import (
    Any,
    Callable,
)

class Decorator:
    _caches = {}

    @classmethod
    def cache(cls, user_function: Callable[..., Any]) -> Callable[..., Any]:
        def decorating_function(*args: Any, **kwargs: Any) -> Any:
            cache = cls._caches.get(user_function.__name__)
            if not cache:
                cache = user_function(*args, **kwargs)
                cls._caches[user_function.__name__] = cache
            return cache

        return decorating_function
